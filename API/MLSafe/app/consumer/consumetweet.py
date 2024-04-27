from confluent_kafka import Consumer, KafkaError
import json
from profanity_check import predict
from elasticsearch import Elasticsearch

def consumetweets():
    conf = {
        'bootstrap.servers': 'localhost:9092',  # Kafka broker address
        'group.id': 'safecheckgroup',        # Consumer group ID
        'auto.offset.reset': 'earliest' ,   
        'enable.auto.commit': 'false' ,
    }
    
    consumer = Consumer(conf)
    
    #Subscribe
    consumer.subscribe(['safecheck'])
    
    try:
        while True:
            #Poll for new messages
            msg = consumer.poll(timeout=1.0)
            
            if msg is None:
                continue
            if msg.error():
                if msg.error().code() == KafkaError._PARTITION_EOF:
                    # end of partition
                    continue
                else:
                    # Handle other errors
                    print(msg.error())
                    break
                    
            #Operations
            data = json.loads(msg.value().decode('utf-8'))
            
            #check if text is safe or not , used cuss_inspect library
            safebool = predict([data.get("text")])
            safe = True if safebool[0] ==0 else False 
                
            
            es = Elasticsearch(['http://localhost:9200'])
            
            doc_id = data.get("status_id")
            print("Doc_id: ",doc_id,"Safe: ",safe)
            
            document = es.get(index='tweetscombined',id=doc_id)
            
            #Add safe feild
            document['_source']['safe'] = safe
            
            #Update the document
            es.update(index='tweetscombined',id=doc_id,body={'doc': document['_source']})
                  
             # Manually commit the Kafka message offset
            consumer.commit(message=msg)
    finally:
        #Close the consumer to release resources
        consumer.close()
    
