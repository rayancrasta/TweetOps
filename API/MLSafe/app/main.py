from fastapi import FastAPI
from .routers import tweetcheck
import asyncio
from .consumer import consumetweet

app = FastAPI()

app.include_router(tweetcheck.router)

async def startkafkaconsumer():
    loop = asyncio.get_event_loop()
    loop.run_in_executor(None,consumetweet.consumetweets)

@app.on_event("startup")
async def startup_event():
    asyncio.create_task(startkafkaconsumer())
    
    
@app.get("/")
def get_root():
    return {"message": "Use Routes"}