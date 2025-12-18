import { MouseEventHandler, useEffect, useState } from 'react';
import { AttachAddon } from '@xterm/addon-attach';
import { FitAddon } from '@xterm/addon-fit';
import { Terminal as XtermTerminal } from '@xterm/xterm';
import { useTranslation } from 'react-i18next';

import '@xterm/xterm/css/xterm.css';

import { Button, Form, Input, Modal } from 'antd';
import { useNavigate } from 'react-router-dom';

import { http } from '@/lib/http.ts';
import { getBaseUrl } from '@/lib/service.ts';
import { Head } from '@/components/head.tsx';

import { validatePicocomParameters } from './validater.ts';

export const Terminal = () => {
  const { t } = useTranslation();
  const navigate = useNavigate();

  const [authed, setAuthState] = useState(true);
  const [user, setUser] = useState('');
  const [passwd, setPasswd] = useState('');
  const [modal, setModal] = useState(false);

  useEffect(() => {
    const terminalEle = document.getElementById('terminal');
    if (!terminalEle) return;

    const terminal = new XtermTerminal({
      cursorBlink: true
    });

    const fitAddon = new FitAddon();
    terminal.loadAddon(fitAddon);
    terminal.open(terminalEle);
    fitAddon.fit();

    if (!authed) {
      return;
    }

    const url = `${getBaseUrl('ws')}/api/vm/terminal`;

    const ws = new WebSocket(url);
    let isPicocomRunning = false;

    ws.onopen = () => {
      const attachAddon = new AttachAddon(ws);
      terminal.loadAddon(attachAddon);

      sendSize();
    };

    const sendSize = () => {
      const windowSize = { rows: terminal.rows, cols: terminal.cols };
      const blob = new Blob([JSON.stringify(windowSize)], { type: 'application/json' });
      ws.send(blob);
      setTimeout(runPicocom, 300);
      setTimeout(runAssistant, 300);
    };

    const runPicocom = () => {
      const urls = window.location.href.split('?');
      if (urls.length < 2) return;

      const searchParams = new URLSearchParams(urls[1]);
      const port = searchParams.get('port');
      const baud = searchParams.get('baud');
      const parity = searchParams.get('parity');
      const flowControl = searchParams.get('flowControl');
      const dataBits = searchParams.get('dataBits');
      const stopBits = searchParams.get('stopBits');
      if (!port || !baud) return;

      if (!validatePicocomParameters({ port, baud, parity, flowControl, dataBits, stopBits })) {
        return;
      }

      ws.send(
        `picocom ${port} --baud ${baud} --parity ${parity} --flow ${flowControl} --databits ${dataBits} --stopbits ${stopBits}\r`
      );
      isPicocomRunning = true;
    };

    const runAssistant = () => {
      const urls = window.location.href.split('?');
      if (urls.length < 2) return;

      const searchParams = new URLSearchParams(urls[1]);
      const type = searchParams.get('type');
      const cmd = searchParams.get('cmd');
      if (type !== 'assistant' || cmd !== 'install') return;

      ws.send(`pip install -r /tmp/requirements.txt\r`);
    };

    const exitPicocom = () => {
      if (ws.readyState === WebSocket.OPEN && isPicocomRunning) {
        ws.send('\x01\x18');
        isPicocomRunning = false;
      }
    };

    const cleanupConnection = () => {
      exitPicocom();
      setTimeout(() => {
        if (ws.readyState === WebSocket.OPEN) {
          ws.close();
        }
      }, 100);
    };

    const resizeScreen = () => {
      fitAddon.fit();
      sendSize();
    };

    const handleBeforeUnload = () => {
      cleanupConnection();
    };

    window.addEventListener('resize', resizeScreen, false);
    window.addEventListener('beforeunload', handleBeforeUnload);

    return () => {
      terminal.dispose();

      cleanupConnection();

      window.removeEventListener('resize', resizeScreen, false);
      window.removeEventListener('beforeunload', handleBeforeUnload);
    };
  }, [authed]);

  const handlerAuthSubmit: MouseEventHandler<HTMLInputElement> = () => {
    http
      .get('/api/vm/terminal', {}, {
        Authorization: `Basic ${btoa(`${user}:${passwd}`)}`
      } as any)
      .then((res) => {
        if (res.code === -1) {
          navigate('/', { replace: true });
        } else {
          setAuthState(true);
          setModal(false);
        }
      });
  };

  return (
    <>
      <Head title={t('head.terminal')} />
      <Modal
        open={modal}
        title="login your ssh"
        onCancel={() => setModal(false)}
        footer={
          <div className="flex w-full items-center justify-center">
            <Button type="primary" onClick={handlerAuthSubmit}>
              submit
            </Button>
          </div>
        }
      >
        <Form className="flex flex-col gap-2">
          <div className="flex flex-col gap-1">
            <span>username</span>
            <Input
              value={user}
              onChange={(val) => setUser(val.target.value)}
              aria-label="username"
            />
          </div>
          <div className="flex flex-col gap-1">
            <span>password</span>
            <Input.Password
              value={passwd}
              onChange={(val) => setPasswd(val.target.value)}
              type="password"
              aria-label="password"
            />
          </div>
        </Form>
      </Modal>

      <div className="h-full w-full overflow-hidden">
        <div id="terminal" className="h-full p-2"></div>
      </div>
    </>
  );
};
