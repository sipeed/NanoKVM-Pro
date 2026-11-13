type CharKey = {
  code: string;
  shift?: boolean;
};

const CHAR_KEYS: Record<string, CharKey> = {
  a: { code: 'KeyA' },
  b: { code: 'KeyB' },
  c: { code: 'KeyC' },
  d: { code: 'KeyD' },
  e: { code: 'KeyE' },
  f: { code: 'KeyF' },
  g: { code: 'KeyG' },
  h: { code: 'KeyH' },
  i: { code: 'KeyI' },
  j: { code: 'KeyJ' },
  k: { code: 'KeyK' },
  l: { code: 'KeyL' },
  m: { code: 'KeyM' },
  n: { code: 'KeyN' },
  o: { code: 'KeyO' },
  p: { code: 'KeyP' },
  q: { code: 'KeyQ' },
  r: { code: 'KeyR' },
  s: { code: 'KeyS' },
  t: { code: 'KeyT' },
  u: { code: 'KeyU' },
  v: { code: 'KeyV' },
  w: { code: 'KeyW' },
  x: { code: 'KeyX' },
  y: { code: 'KeyY' },
  z: { code: 'KeyZ' },
  A: { code: 'KeyA', shift: true },
  B: { code: 'KeyB', shift: true },
  C: { code: 'KeyC', shift: true },
  D: { code: 'KeyD', shift: true },
  E: { code: 'KeyE', shift: true },
  F: { code: 'KeyF', shift: true },
  G: { code: 'KeyG', shift: true },
  H: { code: 'KeyH', shift: true },
  I: { code: 'KeyI', shift: true },
  J: { code: 'KeyJ', shift: true },
  K: { code: 'KeyK', shift: true },
  L: { code: 'KeyL', shift: true },
  M: { code: 'KeyM', shift: true },
  N: { code: 'KeyN', shift: true },
  O: { code: 'KeyO', shift: true },
  P: { code: 'KeyP', shift: true },
  Q: { code: 'KeyQ', shift: true },
  R: { code: 'KeyR', shift: true },
  S: { code: 'KeyS', shift: true },
  T: { code: 'KeyT', shift: true },
  U: { code: 'KeyU', shift: true },
  V: { code: 'KeyV', shift: true },
  W: { code: 'KeyW', shift: true },
  X: { code: 'KeyX', shift: true },
  Y: { code: 'KeyY', shift: true },
  Z: { code: 'KeyZ', shift: true },
  '1': { code: 'Digit1' },
  '2': { code: 'Digit2' },
  '3': { code: 'Digit3' },
  '4': { code: 'Digit4' },
  '5': { code: 'Digit5' },
  '6': { code: 'Digit6' },
  '7': { code: 'Digit7' },
  '8': { code: 'Digit8' },
  '9': { code: 'Digit9' },
  '0': { code: 'Digit0' },
  '!': { code: 'Digit1', shift: true },
  '@': { code: 'Digit2', shift: true },
  '#': { code: 'Digit3', shift: true },
  $: { code: 'Digit4', shift: true },
  '%': { code: 'Digit5', shift: true },
  '^': { code: 'Digit6', shift: true },
  '&': { code: 'Digit7', shift: true },
  '*': { code: 'Digit8', shift: true },
  '(': { code: 'Digit9', shift: true },
  ')': { code: 'Digit0', shift: true },
  ' ': { code: 'Space' },
  '\n': { code: 'Enter' },
  '\t': { code: 'Tab' },
  '-': { code: 'Minus' },
  _: { code: 'Minus', shift: true },
  '=': { code: 'Equal' },
  '+': { code: 'Equal', shift: true },
  '[': { code: 'BracketLeft' },
  '{': { code: 'BracketLeft', shift: true },
  ']': { code: 'BracketRight' },
  '}': { code: 'BracketRight', shift: true },
  '\\': { code: 'Backslash' },
  '|': { code: 'Backslash', shift: true },
  ';': { code: 'Semicolon' },
  ':': { code: 'Semicolon', shift: true },
  "'": { code: 'Quote' },
  '"': { code: 'Quote', shift: true },
  '`': { code: 'Backquote' },
  '~': { code: 'Backquote', shift: true },
  ',': { code: 'Comma' },
  '<': { code: 'Comma', shift: true },
  '.': { code: 'Period' },
  '>': { code: 'Period', shift: true },
  '/': { code: 'Slash' },
  '?': { code: 'Slash', shift: true }
};

const RU_LETTERS: Record<string, string> = {
  ё: '`',
  й: 'q',
  ц: 'w',
  у: 'e',
  к: 'r',
  е: 't',
  н: 'y',
  г: 'u',
  ш: 'i',
  щ: 'o',
  з: 'p',
  х: '[',
  ъ: ']',
  ф: 'a',
  ы: 's',
  в: 'd',
  а: 'f',
  п: 'g',
  р: 'h',
  о: 'j',
  л: 'k',
  д: 'l',
  ж: ';',
  э: "'",
  я: 'z',
  ч: 'x',
  с: 'c',
  м: 'v',
  и: 'b',
  т: 'n',
  ь: 'm',
  б: ',',
  ю: '.'
};

export function charToCodes(char: string): string[] | null {
  const direct = CHAR_KEYS[char];
  if (direct) {
    return direct.shift ? ['ShiftLeft', direct.code] : [direct.code];
  }

  const translated = translateCyrillicChar(char);
  if (!translated) {
    return null;
  }

  const mapped = CHAR_KEYS[translated];
  if (!mapped) {
    return null;
  }
  return mapped.shift ? ['ShiftLeft', mapped.code] : [mapped.code];
}

function translateCyrillicChar(ch: string): string | null {
  if (ch === '№') {
    return '#';
  }
  if (ch === 'Ё') {
    return '~';
  }

  const lower = ch.toLowerCase();
  const translated = RU_LETTERS[lower];
  if (!translated) {
    return null;
  }

  return ch === lower ? translated : translated.toUpperCase();
}
