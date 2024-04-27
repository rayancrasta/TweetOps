from fastapi import APIRouter
from .. import schemas 

router = APIRouter()

@router.get("/safecheck/",status_code=200,response_model=schemas.TweetSafeResponse)
def safecheck(tweetsafe : schemas.TweetSafe):
    resp = schemas.TweetSafeResponse
    resp.tweetid = tweetsafe.tweetid
    resp.content = "blahb la"
    resp.safe = True
    return resp