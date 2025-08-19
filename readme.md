# Catchall Errorpage...
... is a silly project used for my Traefik reverse Proxy in case something goes wrong. Feel free to use it too. Docker image supports ARM and AMD64 CPUs.

`docker pull jmgitde/errorpage:latest`

Port is http 8080

## Supported Themes

- xp.html
- 7.html
- 98.html
- classic.html (fallback, default)

## Custom Themes


To use a custom theme, set the `PAGE_THEME` environment variable to the name of the theme file like (powershell example) $env:PAGE_THEME = "7.html"

### xp.html
<img width="1241" height="1334" alt="image" src="https://github.com/user-attachments/assets/d0950fa9-9653-4f66-94b5-f36d04cfd65c" />

### 7.html
<img width="1381" height="1491" alt="image" src="https://github.com/user-attachments/assets/0f1c6846-6efb-445c-8a8d-7b53664fcd8a" />

### 98.html
<img width="1323" height="1321" alt="image" src="https://github.com/user-attachments/assets/3226161c-c829-4a63-a6f0-aa6d71cd3ea3" />

### classic.html
<img width="1331" height="1080" alt="image" src="https://github.com/user-attachments/assets/209adebb-706f-4852-a0a9-7d1ff8845e4b" />

