import datetime
import logging
import secrets

from telethon import TelegramClient, functions, events, sync, types

logging.basicConfig(level=logging.ERROR)

client = TelegramClient('Maxr1998 Userbot', secrets.API_ID, secrets.API_HASH)
client.start()

chat = next((c for c in client.iter_dialogs() if c.title == "Frankfurt - RR"), None)

if chat == None:
    exit()
print(chat.id)

@client.on(events.NewMessage(chats=chat))
async def handler(event):
    if event.raw_text != None and event.raw_text != "":
        await client.send_message(-1001428954352, event.raw_text)

client.run_until_disconnected()

