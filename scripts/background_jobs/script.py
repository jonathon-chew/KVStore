import redis

r = redis.Redis(
    host="localhost",
    port=6379,
    decode_responses=True,
    protocol=2,
)

# SET
print(r.set("name", "Jon"))

# GET
print(r.get("name"))

# INCR
print(r.incr("counter"))

# GETSET
print(r.getset("counter", "4"))

# STRLEN
print(r.strlen("name"))

# GETDEL
print(r.getdel("name"))

# DELETE
print(r.delete("counter"))
