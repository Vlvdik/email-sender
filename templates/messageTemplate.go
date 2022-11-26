<!doctype html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport"
          content="width=device-width, user-scalable=no, initial-scale=1.0, maximum-scale=1.0, minimum-scale=1.0">
    <meta http-equiv="X-UA-Compatible" content="ie=edge">
    <title>Document</title>
</head>
<body style="background: #280441; text-align: center;">
    <h1 style="color:#d7d7d7;"> Привет, {{.PersonalInfo.Name}} {{.PersonalInfo.LastName}}!</h1>
    <h3 style="color: aliceblue">Я знаю, что твоя родился {{.PersonalInfo.Birthday}} </h3>
    <p>В честь этого мы подготовили тебе подарок. Кликай <strong><a href="Адрес сервера/linkTracker?email={{.Email}}">СЮДА</a></strong> чтобы получить приз!</p>
    <img src="Адрес сервера/pixelTracker?email={{.Email}}"  alt="Bird" width="1px" height="1px">
</body>
</html>
