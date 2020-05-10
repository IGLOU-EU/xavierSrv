# ♿ xavierSrv
*Teeny Weeny status checker.*  

---

**Xavier isn’t a simple guy, but the X-Men leader ! With Cerebro, he can watch all your X-services, check their http status, and when something fails… he uses psychic energy to execute some prepared X-request.**  

## 🛃 Usage
Can be use with none, one or two args.  
- First arg is for HTML output **OR** *-h (print help)*
- Second arg for url list path
**Exemple**
``./xavier.sh /x-mansion/x-services.html /x-mansion/cerebro.list``  

By default, when no one arg is given `${root}/status.html` and `${root}/url.list` are used.

## 🛂 url.list
Each line is a tested url, composed like `STATUS:URL[:PORT]`:  
- Expected HTTP status (if not = fail)
- Service url or ip
- Can be add port like `:8080`  

Except all lines prefixed by `#` **OR **`%`
- `#` is for comments
- `%` is for commands when a check fail
