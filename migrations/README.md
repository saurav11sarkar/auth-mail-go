# UUID migration

নতুন database-এ users table-এর ID:

```sql
id UUID PRIMARY KEY DEFAULT gen_random_uuid()
```

সব up migration ক্রমানুসারে চালাতে হবে। প্রথম migration-এ UUID default আছে। আগে তৈরি database-এ সেই file পরিবর্তন করলে schema নিজে বদলায় না; তাই `20261001000000_users_id_default.up.sql` রাখা হয়েছে। Default আগে থেকেই থাকলেও এই ALTER চালানো যায়। পুরোনো row-এর ID বদলাবে না।

প্রকল্পের migration CLI এবং Make থাকলে, সঠিক DATABASE_URL সেট করে:

```sh
make migrate-up
```

নতুন application চালানোর আগে migration apply করুন। Database reset/drop প্রয়োজন নেই।

Registration request-এ ID দিতে হবে না। Repository INSERT-এ id বাদ দেয় এবং RETURNING id দিয়ে database-generated UUID model-এ নেয়। Response, JWT, profile lookup এবং সম্পর্কের জন্য ID field থাকবে।

এই project-এ বর্তমানে শুধু users table আছে; blog table বা blog create query নেই। ভবিষ্যতে blogs-এর নিজস্ব primary key-তে একই UUID default দেওয়া যায়। তবে author_id-এর মতো foreign key-তে নতুন UUID generate করবেন না—সেখানে existing user ID দিতে হবে।

ID-default down migration পুরোনো application-generated-ID schema-তে ফেরায়। সেটি চালানোর আগে application-কে ID পাঠানোর উপযোগী version-এ rollback করতে হবে।
