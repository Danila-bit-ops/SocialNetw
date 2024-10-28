document.addEventListener('DOMContentLoaded', () => {
    const registerButton = document.getElementById('register-button');
    const loginButton = document.getElementById('login-button');
    const logoutButton = document.getElementById('logout-button'); // Добавлено

    const registerForm = {
        login: document.getElementById('register-login'),
        email: document.getElementById('register-email'),
        password: document.getElementById('register-password')
    };

    const loginForm = {
        email: document.getElementById('login-email'),
        password: document.getElementById('login-password')
    };

    // Функция для сохранения токена в localStorage
    function saveToken(token) {
        localStorage.setItem('token', token);
    }

    // Функция для получения токена из localStorage
    function getToken() {
        return localStorage.getItem('token');
    }

    // Функция для удаления токена
    function removeToken() {
        localStorage.removeItem('token');
    }

    // Функция для проверки аутентификации
    function isAuthenticated() {
        return !!getToken();
    }

    // Функция для отображения новостей
    function showNews() {
        document.body.classList.add('authenticated');
        document.getElementById('auth').style.display = 'none';
        loadNews();
    }

    // Функция для скрытия новостей и отображения форм
    function showAuthForms() {
        document.body.classList.remove('authenticated');
        document.getElementById('auth').style.display = 'block';
        // Очистка списка новостей при выходе
        const latestNewsSection = document.getElementById('последние-новости');
        latestNewsSection.querySelector('ul').innerHTML = '<li class="loader">Загрузка новостей...</li>';
    }

    // Обработчик регистрации
    registerButton.addEventListener('click', (e) => {
        e.preventDefault();
        const login = registerForm.login.value.trim();
        const email = registerForm.email.value.trim();
        const password = registerForm.password.value.trim();

        if (!login || !email || !password) {
            alert('Пожалуйста, заполните все поля для регистрации.');
            return;
        }

        fetch('/api/register', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ login, email, password })
        })
        .then(response => response.json())
        .then(data => {
            if (data.error) {
                alert(`Ошибка: ${data.error}`);
            } else {
                alert('Регистрация прошла успешно. Пожалуйста, войдите в систему.');
                // Очистка полей формы
                registerForm.login.value = '';
                registerForm.email.value = '';
                registerForm.password.value = '';
            }
        })
        .catch(error => {
            console.error('Ошибка при регистрации:', error);
            alert('Произошла ошибка при регистрации. Попробуйте позже.');
        });
    });

    // Обработчик входа
    loginButton.addEventListener('click', (e) => {
        e.preventDefault();
        const email = loginForm.email.value.trim();
        const password = loginForm.password.value.trim();

        if (!email || !password) {
            alert('Пожалуйста, заполните все поля для входа.');
            return;
        }

        fetch('/api/login', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ email, password })
        })
        .then(response => response.json())
        .then(data => {
            if (data.error) {
                alert(`Ошибка: ${data.error}`);
            } else {
                saveToken(data.token);
                showNews();
            }
        })
        .catch(error => {
            console.error('Ошибка при входе:', error);
            alert('Произошла ошибка при входе. Попробуйте позже.');
        });
    });

    // Обработчик выхода из системы
    logoutButton.addEventListener('click', (e) => {
        e.preventDefault();
        removeToken();
        showAuthForms();
        alert('Вы успешно вышли из системы.');
    });

    // Обработчик добавления нового поста
    const addPostButton = document.getElementById('add-post-button');
    const postTitleInput = document.getElementById('post-title');
    const postTextInput = document.getElementById('post-text');

    addPostButton.addEventListener('click', (e) => {
        e.preventDefault();
        const title = postTitleInput.value.trim();
        const post_text = postTextInput.value.trim();

        if (!title || !post_text) {
            alert('Пожалуйста, заполните все поля для добавления поста.');
            return;
        }

        fetch('/api/posts', { // Предполагаемый эндпоинт для добавления поста
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': getToken()
            },
            body: JSON.stringify({ title, post_text })
        })
        .then(response => response.json())
        .then(data => {
            if (data.error) {
                alert(`Ошибка: ${data.error}`);
            } else {
                alert('Пост успешно добавлен.');
                // Очистка полей формы
                postTitleInput.value = '';
                postTextInput.value = '';
                // Перезагрузка списка новостей
                loadNews();
            }
        })
        .catch(error => {
            console.error('Ошибка при добавлении поста:', error);
            alert('Произошла ошибка при добавлении поста. Попробуйте позже.');
        });
    });

    // Функция загрузки новостей
    function loadNews() {
        fetch('/api/news', {
            headers: {
                'Authorization': getToken()
            }
        })
        .then(response => {
            if (!response.ok) {
                throw new Error(`Ошибка сети: ${response.status}`);
            }
            return response.json();
        })
        .then(data => {
            console.log('Загруженные новости:', data); // Для отладки
            const latestNewsSection = document.getElementById('последние-новости');
            const newsList = latestNewsSection.querySelector('ul');

            // Очистка существующих новостей
            newsList.innerHTML = '';

            if (data.length === 0) {
                const noPostsItem = document.createElement('li');
                noPostsItem.textContent = 'Новостей пока нет.';
                newsList.appendChild(noPostsItem);
                return;
            }

            data.forEach(post => {
                const listItem = document.createElement('li');

                const title = document.createElement('h3');
                title.textContent = post.title || 'Без заголовка';
                listItem.appendChild(title);

                const description = document.createElement('p');
                description.textContent = post.post_text || 'Нет описания';
                listItem.appendChild(description);

                const readMoreLink = document.createElement('a');
                readMoreLink.href = `/news/${post.post_id}`; // Замените на реальную ссылку
                readMoreLink.textContent = 'Подробнее &rarr;';
                listItem.appendChild(readMoreLink);

                newsList.appendChild(listItem);
            });
        })
        .catch(error => {
            console.error('Ошибка при загрузке новостей:', error);
            const latestNewsSection = document.getElementById('последние-новости');
            latestNewsSection.querySelector('ul').innerHTML = '<li>Не удалось загрузить новости. Попробуйте позже.</li>';
        });
    }

    // Проверка аутентификации при загрузке страницы
    if (isAuthenticated()) {
        showNews();
    } else {
        showAuthForms();
    }
});
