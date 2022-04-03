# Expected psuedo code which should be described using DAG

## open top (and login)

```
open https://theminjoo.kr/
if ("div.hds_right > a.hd_login").contains("로그인 해주세요")
	click("div.hds_right > a.hd_login")
fi
wait("form#frm_login", 10s)
input("form#frm_login > div  > div > input#login_id", $myId)
input("form#frm_login > div  > div > input#login_pw", $myPwd)
click("form#frm_login > div  > div > label > input#login_save")
click("form#frm_login > div  > div > button.submit")
```