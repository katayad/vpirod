import mincemeat

client = mincemeat.Client()
client.password	= "changeme"
client.conn("localhost", mincemeat.DEFAULT_PORT)
