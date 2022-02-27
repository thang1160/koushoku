import "./styles/main.less";

interface PageState {
  isViewing?: boolean;
  isPreloading?: boolean;
  isPreloaded?: boolean;
}

const pages: PageState[] = [];

let id: string;
let slug: string;

const initPreload = () => {
  const currentIndex = pages.findIndex(e => e.isViewing);

  for (let i = 0; i < pages.length; i++) {
    const page = pages[i];
    if (
      page.isPreloading ||
      page.isPreloaded ||
      i < Math.max(0, currentIndex - 3) ||
      i > Math.max(0, currentIndex + 3)
    ) {
      continue;
    }
    page.isPreloading = true;

    const img = document.createElement("img");
    img.src = `/data/${id}/${i + 1}.jpg`;

    const onComplete = () => {
      page.isPreloading = false;
      page.isPreloaded = true;
    };
    const onFailed = () => (page.isPreloading = false);

    if (img.complete) {
      onComplete();
    } else {
      img.addEventListener("load", onComplete);
      img.addEventListener("error", onFailed);
    }
  }
};

const initReader = () => {
  const reader = document.getElementById("reader");
  if (!reader) return;

  const arr = window.location.pathname.split("/");
  [, , id, slug] = arr;

  const total = Number(document.querySelector(".total").textContent);
  const currentSpans = document.querySelectorAll(".current");
  let current = Number(currentSpans[0].textContent);

  const first = document.querySelectorAll(".first") as NodeListOf<HTMLAnchorElement>;
  const last = document.querySelectorAll(".last") as NodeListOf<HTMLAnchorElement>;
  const prev = document.querySelectorAll(".prev") as NodeListOf<HTMLAnchorElement>;
  const next = document.querySelectorAll(".next") as NodeListOf<HTMLAnchorElement>;

  const page = reader.querySelector(".page a") as HTMLAnchorElement;
  let img = page.querySelector("img");

  const mutex = { current: false };
  const rect = page.getBoundingClientRect();

  pages.push(...Array.from({ length: total }, (_, i) => ({ isViewing: i + 1 === current })));

  const changePage = (n: number) => {
    if (mutex.current) return;
    mutex.current = true;
    try {
      current = n;

      pages.find(p => p.isViewing).isViewing = false;
      pages[current - 1].isViewing = true;

      page.href = `/archive/${id}/${slug}/${current}`;

      const newImg = document.createElement("img");
      newImg.src = `/data/${id}/${current}.jpg`;
      img.replaceWith(newImg);
      img = newImg;

      currentSpans.forEach(span => (span.textContent = current.toString()));
      window.history.replaceState(null, "", `/archive/${id}/${slug}/${current}`);

      prev.forEach(e => {
        e.href = `/archive/${id}/${slug}/${current - 1}`;
      });

      next.forEach(e => {
        e.href = `/archive/${id}/${slug}/${current + 1}`;
      });
    } finally {
      mutex.current = false;
    }
  };

  first.forEach(e => {
    e.addEventListener("click", ev => {
      ev.preventDefault();
      ev.stopPropagation();
      ev.stopImmediatePropagation();

      if (current === 1) return;
      changePage(1);
    });
  });

  last.forEach(e => {
    e.addEventListener("click", ev => {
      ev.preventDefault();
      ev.stopPropagation();
      ev.stopImmediatePropagation();

      if (current === total) return;
      changePage(total);
    });
  });

  prev.forEach(e => {
    e.addEventListener("click", ev => {
      ev.preventDefault();
      ev.stopPropagation();
      ev.stopImmediatePropagation();

      if (current === 1) return;
      changePage(current - 1);
    });
  });

  next.forEach(e => {
    e.addEventListener("click", ev => {
      ev.preventDefault();
      ev.stopPropagation();
      ev.stopImmediatePropagation();

      if (current === total) return;
      changePage(current + 1);
    });
  });

  page.addEventListener("click", (ev: MouseEvent) => {
    ev.preventDefault();
    ev.stopPropagation();
    ev.stopImmediatePropagation();

    const target = ev.target as HTMLAnchorElement;
    const isPrev = ev.screenX - rect.x <= target.clientWidth / 2;

    if (isPrev) {
      if (current === 1) return;
      changePage(current - 1);
    } else {
      if (current === total) return;
      changePage(current + 1);
    }
  });

  let interval = 0;
  const scrollIntoView = () => {
    clearInterval(interval);
    interval = window.setInterval(() => {
      window.requestAnimationFrame(() => {
        if (!img.naturalHeight) return;
        clearInterval(interval);
        page.scrollIntoView({ block: "start", inline: "start" });
      });
    }, 10);
  };
  scrollIntoView();

  const observer = new MutationObserver(() => {
    scrollIntoView();
    initPreload();
  });
  observer.observe(page, { childList: true, attributes: true });
  initPreload();

  window.addEventListener("keydown", ev => {
    if (ev.code === "ArrowLeft" || ev.code === "KeyA" || ev.code === "KeyH") {
      if (current === 1) return;
      changePage(current - 1);
    } else if (ev.code === "ArrowRight" || ev.code === "KeyD" || ev.code === "KeyL") {
      if (current === total) return;
      changePage(current + 1);
    }
  });
};

const init = () => {
  if ("serviceWorker" in navigator) {
    navigator.serviceWorker.register("/serviceWorker.js");
  }
  initReader();
};

if (document.readyState === "complete") {
  init();
} else {
  document.addEventListener("DOMContentLoaded", init);
}
