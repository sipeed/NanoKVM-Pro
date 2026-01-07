import { atom } from 'jotai';

// is the settings modal opened
export const isSettingsOpenAtom = atom(false);

// current setting tab
export const settingTabAtom = atom('about');

// web title
export const webTitleAtom = atom('');

// menu bar disabled items
export const menuDisabledItemsAtom = atom<string[]>([]);

// track how many submenus are currently open
export const submenuOpenCountAtom = atom(0);
