function setupAppLayoutExamples() {
  var aside = document.querySelector('.l-aside');
  var navigation = document.querySelector('.l-navigation');

  var menuToggle = document.querySelector('.js-menu-toggle');
  var menuClose = document.querySelector('.js-menu-close');
  var menuPin = document.querySelector('.js-menu-pin');
  var asideOpen = document.querySelector('.js-aside-open');
  var asideClose = document.querySelector('.js-aside-close');
  var asideResize = document.querySelectorAll('.js-aside-resize');
  var asidePin = document.querySelector('.js-aside-pin');

  if (menuToggle) {
    menuToggle.addEventListener('click', function () {
      navigation.classList.toggle('is-collapsed');
    });
  }

  if (menuClose) {
    menuClose.addEventListener('click', function (e) {
      navigation.classList.add('is-collapsed');
      document.activeElement.blur();
    });
  }

  if (asideOpen) {
    asideOpen.addEventListener('click', function () {
      aside.classList.remove('is-collapsed');
    });
  }

  if (asideClose) {
    asideClose.addEventListener('click', function () {
      aside.classList.add('is-collapsed');
    });
  }

  [].slice.call(asideResize).forEach(function (button) {
    button.addEventListener('click', function () {
      button.dataset.resizeClass;
      var panel = document.getElementById(button.getAttribute('aria-controls'));
      if (panel) {
        panel.classList.remove('is-narrow');
        panel.classList.remove('is-medium');
        panel.classList.remove('is-wide');
        if (button.dataset.resizeClass) {
          panel.classList.add(button.dataset.resizeClass);
        }
      }
    });
  });

  if (menuPin) {
    menuPin.addEventListener('click', function () {
      navigation.classList.toggle('is-pinned');
      if (navigation.classList.contains('is-pinned')) {
        menuPin.querySelector('i').classList.add('p-icon--close');
        menuPin.querySelector('i').classList.remove('p-icon--pin');
      } else {
        menuPin.querySelector('i').classList.add('p-icon--pin');
        menuPin.querySelector('i').classList.remove('p-icon--close');
      }
      document.activeElement.blur();
    });
  }

  if (asidePin) {
    asidePin.addEventListener('click', function () {
      aside.classList.toggle('is-pinned');
      if (aside.classList.contains('is-pinned')) {
        asidePin.innerText = 'Unpin';
      } else {
        asidePin.innerText = 'Pin';
      }
    });
  }
}

function submitModal(event) {
  event.preventDefault();
  const form = document.getElementById('modal-form');
  const formData = new FormData(form);
  const method = form.attributes['data-method'].value;
  if (method === 'PUT') {
    const id = form.attributes['data-id'].value;
    formData.append('id', id);
    fetch(form.action, {
      method: 'PUT',
      body: formData,
    })
      .then(response => {
        if (response.ok) {
          console.log('Repository updated successfully');
          const url = new URL(form.action);
          const pathParts = url.pathname.split('/');
          pathParts.pop();
          url.pathname = pathParts.join('/');
          window.location.href = url.href;
        } else {
          console.error('Failed to update repository');
        }
      })
      .catch(error => {
        console.error('Network error:', error);
      });
  }
  form.submit().then(response => {
    if (response.ok) {
      console.log('Form submitted successfully');
      window.location.reload();
    } else {
      console.error('Failed to submit form');
    }
  })
    .catch(error => {
      console.error('Network error:', error);
    });
}

function deleteObject(id) {
  fetch(window.location.href, {
    method: 'DELETE', body: {
      id: id
    },
  })
    .then(response => {
      if (response.ok) {
        console.log('Object deleted successfully');
        window.location.href = response.headers.get('Location');
      } else {
        console.error('Failed to delete object');
      }
    })
    .catch(error => {
      console.error('Network error:', error);
    });
}

setupAppLayoutExamples();
