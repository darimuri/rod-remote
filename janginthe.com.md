implement using https://github.com/deliveryhero/pipeline
explained at https://golangexample.com/a-library-to-help-you-create-pipelines-in-golang/


https://janginthe.com/

# 로그인 블록

* 로그아웃 상태
```
document.querySelector("#header > div.inner > div > div > ul.xans-element-.xans-layout.xans-layout-statelogoff.right.off > li:nth-child(1) > a") 로그아웃
```
* 로그인 상태
```
document.querySelector("#header > div.inner > div > div > ul.xans-element-.xans-layout.xans-layout-statelogon.right.in > li:nth-child(1) > a") 로그인
```

## 로그인 화면
https://janginthe.com/member/login.html
* 아이디
```
document.querySelector("#member_id")
```
* 비밀번호
```
document.querySelector("#member_passwd")
```
* 로그인 버튼
```
document.querySelector("#member_form_1658865576 > div > div > fieldset > a")
```

# 공지 팝업
* 오늘 하루 체크박스
```
document.querySelector("input#popup_close_check")
```
* 닫기 버튼
```
document.querySelector("#popup_close_btn")
```


# 제품

## 품절
[의정부 장인한과 못난이 약과 (파지약과)](https://janginthe.com/product/%EC%9D%98%EC%A0%95%EB%B6%80-%EC%9E%A5%EC%9D%B8%ED%95%9C%EA%B3%BC-%EB%AA%BB%EB%82%9C%EC%9D%B4-%EC%95%BD%EA%B3%BC-%ED%8C%8C%EC%A7%80%EC%95%BD%EA%B3%BC/260/category/28/display/1/)
* 품절버튼
```
document.querySelector("#contents > div.xans-element-.xans-product.xans-product-detail > div.detailArea > div.infoArea > div.xans-element-.xans-product.xans-product-action > div.ec-base-button > span")
```

## 구매 가능
[장인,더 약과빵](https://janginthe.com/product/%EC%9E%A5%EC%9D%B8%EB%8D%94-%EC%95%BD%EA%B3%BC%EB%B9%B5/258/category/28/display/1/)
* 구매버튼
```
document.querySelector("#contents > div.xans-element-.xans-product.xans-product-detail > div.detailArea > div.infoArea > div.xans-element-.xans-product.xans-product-action > div.ec-base-button > a.first")
```

### 상세
* 수량 입력
```
document.querySelector("#quantity")
```

## 주문/결제
[주문/결제](https://janginthe.com/order/orderform.html?basket_type=A0000&delvtype=A)
* 결제버튼
```
document.querySelector("#btn_payment")
```
* 결제정보 확인버튼
```
document.querySelector("#stdPaymentConfirmView > div.pgmBtnArea > a > img")
```