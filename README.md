<div align="center">
    <img src=".media/kirby.png" height="300" />
    <h1>~ Kirby ~</h1>
    <strong>
        A Discord Data collection bot that respects users privacy.
    </strong><br><br>
    <a href="https://dc.zekro.de"><img height="28" src="https://img.shields.io/discord/307084334198816769.svg?style=for-the-badge&logo=discord" /></a>&nbsp;
    <img height="28" src="https://forthebadge.com/images/badges/60-percent-of-the-time-works-every-time.svg" />&nbsp;
    <img height="28" src="https://forthebadge.com/images/badges/built-with-grammas-recipe.svg">
<br>
</div>

---

## Introduction

Kirby is a Discord Bot designed to collect general user behaviour data of your Discord guilds members. The collected data contains **no** information about specific user ID's, names, nicknames or message content to respect the users privacy.

## Collected Data

### Which data will be collected?

- **Messages**  
  Every time a message was sent to a channel visible to Kirby, the following properties will be saved:  
  - `timestamp` when the message was sent
  - `channelID` of the text channel the message was sent into
  - `contentLength` of the messages content
  - `bot` identifies wether the sender was a bot
  - `roleIDs` containing the IDs of the roles of the sender
  - `mentions` identifies the number of mentions in the message
  - `mentionedRoles` cintains the IDs of mentioned roles
  - `attachment` identifies if the message contains any attachments
  - `guildID` to identify the guild where the message was sent

- **Status Updates**  
  Every time a user updates their online status on a guild, the following properties will be recorded:  
  - `roleIDs` containing the IDs of the roles of a user
  - `bot` identifies wether the user was a bot
  - `timestamp` when the status update occured
  - `oldStatus` the status before the update
  - `newStatus` the status after the update
  - `guildID` to identify the guild of the user

- **Member Changes**  
  Every time a user joins or leaves a guild, the following data will be collected:  
  - `roleIDs` containing the IDs of the roles the user had *(empty when user joins ofc)*
  - `timestamp` when the event occured
  - `bot` identifies wether the user is a bot
  - `event` numeral value if the event was a join or a quit event
  - `guildID` to identify the guild where the event occured

- **Reactions**  
  Each time a reaction was added to a message, the following properties will be collected:  
  - `channelID` of the channel of the message which was reacted to
  - `roleIDs` contains the IDS of the users roles who reacted
  - `emoji` contains the UTF-8 name of the emoji
  - `emojiID` is the ID of the emoji if it's a custom emoji
  - `contentLen` is the length of the content of the message reacted to
  - `bot` identifies wther the user reacted was a bot
  - `timestamp` when the event occured
  - `guildID` identifies the guild where the event occured

- **Voice Channel Events**  
  Each time a user joins, leaves or moved between voice channels, the follwing data will be obtained:  
  - `channelID` of the voice channel (joined channel when moved)
  - `roleIDs` containing the IDs of the roles of the user
  - `event` is the numeral identifier which event occured
  - `timestamp` when the event occured
  - `guildID` to identify the guild where the event occured

- **Member Count Stats**  
  Every 30 minutes, the bot collects the count of users per guild for each status and also another times for each different role.

## Which data will NOT be collected?

- **NO** account specific data like userID, name, nichname, tag or discriminator will be obtained
- **NO** specific message content data will be obtained
- **NO** specific behaviour of specific users will be tracked

Further details about the project can be found in the following video (language is in german!):

[![](https://img.youtube.com/vi/mvTeEEeb0jM/0.jpg)](https://www.youtube.com/watch?v=mvTeEEeb0jM)

