from pydantic import BaseModel

class TweetSafe(BaseModel):
    tweetid : int
    
class TweetSafeResponse(TweetSafe):
    content: str
    safe : bool
    