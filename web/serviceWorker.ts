import { CacheableResponsePlugin } from "workbox-cacheable-response";
import { ExpirationPlugin } from "workbox-expiration";
import { registerRoute } from "workbox-routing";
import { CacheFirst, StaleWhileRevalidate } from "workbox-strategies";

registerRoute(
  ({ request }) => request.mode === "navigate",
  new StaleWhileRevalidate({
    cacheName: "pages",
    plugins: [new CacheableResponsePlugin({ statuses: [200] }), new ExpirationPlugin({ maxAgeSeconds: 3600 })]
  })
);

registerRoute(
  ({ request, url }) =>
    request.destination === "style" ||
    request.destination === "script" ||
    request.destination === "worker" ||
    url.pathname.startsWith("/fonts/"),
  new StaleWhileRevalidate({
    cacheName: "assets",
    plugins: [new CacheableResponsePlugin({ statuses: [200] })]
  })
);

registerRoute(
  ({ url }) => url.pathname.startsWith("/data/"),
  new CacheFirst({
    cacheName: "data",
    plugins: [
      new CacheableResponsePlugin({ statuses: [200] }),
      new ExpirationPlugin({
        maxEntries: 1024,
        maxAgeSeconds: 86400,
        purgeOnQuotaError: true
      })
    ]
  })
);
