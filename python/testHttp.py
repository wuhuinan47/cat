from aiohttp import web
import json

async def handle(request):
    name = request.match_info.get('name', "Anonymous")
    text = "Hello, " + name

    if name=="getJson":
        data = {'test': '测试'}
        return web.json_response(data, dumps=json.dumps)

    return web.Response(text=text)



async def wshandle(request):
    ws = web.WebSocketResponse()
    await ws.prepare(request)

    async for msg in ws:
        if msg.type == web.WSMsgType.text:
            await ws.send_str("Hello, {}".format(msg.data))
        elif msg.type == web.WSMsgType.binary:
            await ws.send_bytes(msg.data)
        elif msg.type == web.WSMsgType.close:
            break

    return ws


app = web.Application()
app.add_routes([web.get('/', handle),
                web.get('/echo', wshandle),
                web.get('/{name}', handle)])

if __name__ == '__main__':
    web.run_app(app)