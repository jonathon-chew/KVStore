import redis

r = redis.Redis(
    host="localhost",
    port=6379,
    decode_responses=True,
    protocol=2,
)

# SET
try:
	print("SET: ", r.set("name", "Jon"))
except Exception as e:
	print("[ERROR]: ", e)

# GET
try:
	print("GET: ", r.get("name"))
except Exception as e:
	print("[ERROR]: ", e)

# INCR
try:
	# print("GET: ", r.get("counter"))
	print("INCR: ", r.incr("counter"))
except Exception as e:
	print("[ERROR]: ", e)

# GETSET
try:
	print("GETSET: ", r.getset("name", "Not Jon"))
except Exception as e:
	print("[ERROR]: ", e)

# STRLEN
try:
	print("STRLEN: ", r.strlen("name"))
except Exception as e:
	print("[ERROR]: ", e)

# GETDEL
try:
	print("GETDEL: ", r.getdel("name"))
except Exception as e:
	print("[ERROR]: ", e)

# DELETE
try:
	print("SET: ", r.set("delete_me", "True"))
	print("DELETE: ", r.delete("delete_me"))
except Exception as e:
	print("[ERROR]: ", e)
