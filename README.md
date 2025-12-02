## 目的

デッドロックを起こすパターンの具体例を書いておくことで、忘れたら戻って来れるようにする

## 現時点での認識

真に排他なロックは排他レコードロック。

排他ロックは基本は、ロックしたリソースを他のトランザクションに触らせない認識。

が、排他ギャップロックは範囲が被っても取れる

### supremumを共に取ってしまうパターン

（その結果、supremumをとってしまいその後insert intention lock waitinになり最後まで処理が走らないから排他ロックが解除されず、デッドロックになるパターンとか）

一方、共有ロックも複数とって良いが、同じ範囲に対して共有ロックをとり、それぞれのトランザクションがupdateをしようとするとデッドロックになるパターンもある。

（insert,update,deleteは基本排他ロックを取ろうとする）

## 外部キー制約で子供と親を同時に更新した際に起きるパターン

また、外部キー制約でもデッドロックが起きるパターンがある。

子を追加したついでに親も更新するパターン。

parent,childがいるとして、

TX1 -> child1 insert with parent 1 (parent 1にSロック)
TX2 -> child2 insert with parent 1（parent 1にSロック）
TX1 -> parent1 update (ここでXロックに昇格しようとするがTX2がSロックを持つ)
TX2 -> parent1 update（Xロックを要求するがすでにTX1が昇格待ちなのでTX2が終わらない。つまりSロックが手放されない）

## tips

#### 現在のロック状態を可視化するクエリ

```sql
SELECT 
    ENGINE_TRANSACTION_ID as TRX_ID,
    OBJECT_NAME as TABLE_NAME,
    INDEX_NAME,
    LOCK_TYPE,
    LOCK_MODE,
    LOCK_STATUS,
    LOCK_DATA 
FROM performance_schema.data_locks;
```

#### デッドロックログを見る

```sql
SHOW ENGINE INNODB STATUS\G;
```