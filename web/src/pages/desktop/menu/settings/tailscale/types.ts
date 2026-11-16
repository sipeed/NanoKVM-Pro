export type State = 'notInstall' | 'notRunning' | 'notLogin' | 'stopped' | 'running';

export type Status = {
  state: State;
  name: string;
  ip: string;
  account: string;
  serve: boolean;
  serveUrl?: string;
};
