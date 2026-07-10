const DEFAULT_NOTIFICATION_URL = '/';

self.addEventListener('install', (event) => {
  self.skipWaiting();
});

self.addEventListener('activate', (event) => {
  event.waitUntil(self.clients.claim());
});

self.addEventListener('push', (event) => {
  let payload = {};
  if (event.data) {
    try {
      payload = event.data.json();
    } catch (err) {
      payload = { body: event.data.text() };
    }
  }

  const title = payload.title || 'Work order baru masuk';
  const options = {
    body: payload.body || 'Ada work order baru.',
    icon: '/static/public/itlogo.png',
    badge: '/static/public/itlogo.png',
    tag: payload.workOrderId ? `work-order-${payload.workOrderId}` : 'work-order-incoming',
    renotify: true,
    data: {
      url: payload.url || DEFAULT_NOTIFICATION_URL,
      workOrderId: payload.workOrderId || null,
    },
  };

  event.waitUntil(self.registration.showNotification(title, options));
});

self.addEventListener('notificationclick', (event) => {
  event.notification.close();

  const targetUrl = (event.notification.data && event.notification.data.url) || DEFAULT_NOTIFICATION_URL;
  const absoluteUrl = new URL(targetUrl, self.location.origin).href;

  event.waitUntil(
    self.clients.matchAll({ type: 'window', includeUncontrolled: true }).then((clientList) => {
      for (const client of clientList) {
        if (client.url.startsWith(self.location.origin) && 'focus' in client) {
          client.focus();
          if ('navigate' in client && client.url !== absoluteUrl) {
            return client.navigate(absoluteUrl);
          }
          return client;
        }
      }

      if (self.clients.openWindow) {
        return self.clients.openWindow(absoluteUrl);
      }
      return undefined;
    })
  );
});
