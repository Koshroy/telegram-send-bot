telegram-send-bot
=================

Send messages from a named pipe to Telegram. This bot was originally
designed to be used to receive messages over NNCP and send those
received messages to Telegram, however the bot can use any named
pipe. Usage is as follows:

```
telegram-send-bot <fifo>
```

Make sure `TG_BOT_SECRET` is set to the value of your Telegram bot
secret. This bot reads the named pipe line-by-line. Each line is
parsed simply:

```
<chatId> <message contents>
```

The chat ID is the Telegram chat ID, and the message contents are the
actual message to be sent.

## Using with NNCP
This bot requires a named pipe to exist which it can read from to send
data. `fifo-recv.sh` is a simple script which appends output onto a
named pipe and can be used as an exec handle in NNCP to append onto a
named pipe, which the bot can then use to send a Telegram message.
