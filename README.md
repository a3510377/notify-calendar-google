# GCalReminder 日曆提醒系統

GCalReminder 日曆提醒系統，支援 line，Discord

## 運行

### 直接執行

```bash
go run .
```

### Docker

```bash
docker run --name notify-calender-google -d \
  --restart unless-stopped \
  -v /path/to/data:/app/data \
  ghcr.io/a3510377/notify-calendar-google:latest
```

## 設定

第一次啟動時將會生成 `config.yaml` 檔案:

```yaml
CALENDAR_ID: <請輸入自己的日曆編號>
discord:
  enable: false
  TOKEN: ""
  webhook: []
  channel_ids: []
line:
  enable: false
  tokens: []
```

### `CALENDAR_ID` 日曆 ID

該值可通過網址獲取，比如以下範例網址中: `https://calendar.google.com/calendar/u/0/embed?src=c_nbtiskrng1pkrcj168db62l4hg@group.calendar.google.com&ctz=Asia/Taipei` 可以看到 `.../embed?src=c_nbtiskrng1pkrcj168db62l4hg@group.calendar.google.com&ctz=...` 其中將 `?src=` 後方文字複製到 `&ctz=` 前即可。

### discord

Discord 是一款專為社群設計的免費網路即時通話軟體與數位發行平台，可在上方進行通訊/聊天/圖片等傳送等等。

#### `enable` 啟用 discord 提醒

若要關閉 discord 提醒，可將值改成 `false` 若要開啟則改成 `true`

#### `TOKEN` 機器人權杖

若要使用機器人進行發送，需要提供機器人權杖，可於 [#discord/channel_ids](#channel_ids-需發送頻道-id) 添加目標頻道及移除目標頻道

#### `channel_ids` 需發送頻道 ID

**初始值:**

```yaml
channel_ids: []
```

可依您自行添加或刪除頻道 ID。

**範例:**

```yaml
channel_ids:
  - 1008911713875267596
  - 1008911713875267595
  # - 可添加其它更多 ID
```

#### `webhook` Webhook 網址

**初始值:**

```yaml
webhook: []
```

可依您自行添加或刪除 Webhook網址。

**範例:**

```yaml
webhook:
  - https://discord.com/api/webhooks/1008911713875267596/3tM8c2O8Aqa0bkLNzILV3py-TA5RdY3Xy3aG7EkE-iXavEvmO7QL3A15zWkbbd8DAaUH
  - https://discord.com/api/webhooks/1008911713875267595/eLFCf1Gux1SSTUBOteEJNvHBpGaIm8WtcGyDL8gdoZGSskAIjExs01ygU7VBw-NBaAaZ
  # - 可添加其它更多 Webhook 網址
```

### `line` notify

#### `enable` 啟用 line notify 提醒

若要關閉 line notify 提醒，可將值改成 `false` 若要開啟則改成 `true`

#### `tokens` line notify 權杖

若要使用 line notify 進行發送，需要提供 line notify 權杖，可至 [Line Notify](https://notify-bot.line.me/) 申請。

## License

[MIT License](LICENSE)
