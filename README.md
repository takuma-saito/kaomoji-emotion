顔文字の感情抽出
================

顔文字の感情抽出が出来ます。
正解の顔文字集合を学習させて、顔文字の感情を予測できるようにしました。
識別器は単層のパーセプトロンで実装しています。

*使い方:*
main() 関数の StartServer だけ実行させて下さい。
学習には 10秒 ~ 30 秒程時間が掛かるので立ち上がりは遅いです。

    $ go run main.go
    listen 0:0:0:0:8000

    $ nc localhost 8000
    If you want to exit, type 'exit'.
    > (・∀・)
    戸惑い
    > (-_-;)
    無関心
    > (´・ω・｀)
    好感
    > (T_T)
    悲しみ
    > ^^;
    嘲笑
    > m(__)m
    おねだり
    > (^^ゞ
    あいさつ

上記のように顔文字の感情を推定してくれます。
辞める時は exit して下さい。

一応 Licence とか
-----------------

GPL で配布します。



