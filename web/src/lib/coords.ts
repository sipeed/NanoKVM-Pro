export function clientToScreenPercent(
  clientX: number,
  clientY: number
): { x: number; y: number } | null {
  const screen = document.getElementById('screen') as HTMLVideoElement | HTMLImageElement | null;
  if (!screen) {
    return null;
  }

  const { x, y } = getCorrectedCoords(screen, clientX, clientY);
  return {
    x: roundPercent(x * 100),
    y: roundPercent(y * 100)
  };
}

export function percentToHid(percent: number): number {
  const n = Math.max(0, Math.min(100, percent)) / 100;
  return Math.floor(0x7fff * n) + 0x0001;
}

function getCorrectedCoords(
  screen: HTMLVideoElement | HTMLImageElement,
  clientX: number,
  clientY: number
) {
  const rect = screen.getBoundingClientRect();

  const mediaWidth = getMediaWidth(screen);
  const mediaHeight = getMediaHeight(screen);

  if (!mediaWidth || !mediaHeight) {
    return {
      x: clamp01((clientX - rect.left) / rect.width),
      y: clamp01((clientY - rect.top) / rect.height)
    };
  }

  const videoRatio = mediaWidth / mediaHeight;
  const elementRatio = rect.width / rect.height;

  let renderedWidth = rect.width;
  let renderedHeight = rect.height;
  let offsetX = 0;
  let offsetY = 0;

  if (videoRatio > elementRatio) {
    renderedHeight = rect.width / videoRatio;
    offsetY = (rect.height - renderedHeight) / 2;
  } else {
    renderedWidth = rect.height * videoRatio;
    offsetX = (rect.width - renderedWidth) / 2;
  }

  return {
    x: clamp01((clientX - rect.left - offsetX) / renderedWidth),
    y: clamp01((clientY - rect.top - offsetY) / renderedHeight)
  };
}

function getMediaWidth(screen: HTMLVideoElement | HTMLImageElement): number {
  if ('videoWidth' in screen && screen.videoWidth) {
    return screen.videoWidth;
  }
  if ('naturalWidth' in screen && screen.naturalWidth) {
    return screen.naturalWidth;
  }
  return 0;
}

function getMediaHeight(screen: HTMLVideoElement | HTMLImageElement): number {
  if ('videoHeight' in screen && screen.videoHeight) {
    return screen.videoHeight;
  }
  if ('naturalHeight' in screen && screen.naturalHeight) {
    return screen.naturalHeight;
  }
  return 0;
}

function clamp01(value: number): number {
  if (!Number.isFinite(value)) {
    return 0;
  }
  return Math.max(0, Math.min(1, value));
}

function roundPercent(value: number): number {
  return Math.round(clamp01(value / 100) * 10000) / 100;
}
