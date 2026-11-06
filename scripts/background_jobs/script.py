import redis

r = redis.Redis(
    host="localhost",
    port=6379,
    decode_responses=True,
    protocol=2,
)

# SET
print("SET: ", r.set("name", "Jon"))

# GET
print("GET: ", r.get("name"))

# INCR
print("INCR: ", r.incr("counter"))

# GETSET
print("GETSET: ", r.getset("counter", "4"))

# STRLEN
print("STRLEN: ", r.strlen("name"))

# GETDEL
print("GETDEL: ", r.getdel("name"))

# DELETE
print("DELETE: ", r.delete("counter"))
