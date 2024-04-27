from fastapi import APIRouter
from .. import schemas 
from elasticsearch import Elasticsearch

router = APIRouter()

@router.get("/safecheck/",status_code=200,response_model=schemas.TweetSafeResponse)
def safecheck(tweetsafe : schemas.TweetSafe):
    
    es = Elasticsearch(['http://localhost:9200'])
            
    doc_id = tweetsafe.tweetid
    document = es.get(index='tweetscombined',id=doc_id)
              
    resp = schemas.TweetSafeResponse
    resp.tweetid = tweetsafe.tweetid
    resp.content = document['_source']['text']
    resp.safe = document['_source']['safe']
    return resp