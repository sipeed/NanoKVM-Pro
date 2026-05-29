import { useEffect, useRef } from 'react';
import clsx from 'clsx';
import { useAtomValue } from 'jotai';

import * as api from '@/api/stream.ts';
import { mouseStyleAtom } from '@/jotai/mouse';
import { videoParametersAtom } from '@/jotai/screen.ts';

import DirectWorker from './direct.worker.ts?worker';

export const H265Direct = () => {
  const videoParameters = useAtomValue(videoParametersAtom);
  const mouseStyle = useAtomValue(mouseStyleAtom);

  const canvasRef = useRef<any>();
  const workerRef = useRef<Worker | null>(null);

  useEffect(() => {
    if (!canvasRef.current) {
      return;
    }

    const worker = new DirectWorker();
    workerRef.current = worker;

    const offscreen = canvasRef.current.transferControlToOffscreen();
    worker.postMessage({ type: 'init_h265', canvas: offscreen }, [offscreen]);

    const ws = api.directH265();
    ws.binaryType = 'arraybuffer';

    ws.onmessage = (event) => {
      try {
        worker.postMessage({ type: 'ws_message', data: event.data }, [event.data]);
      } catch (error) {
        console.error('Error processing WebSocket message:', error);
      }
    };

    ws.onerror = () => {
      worker.postMessage({ type: 'error' });
    };

    ws.onclose = () => {
      worker.postMessage({ type: 'close' });
    };

    return () => {
      if (ws.readyState === 1) {
        ws.close();
      }
      worker.terminate();
    };
  }, []);

  return (
    <div className="flex h-screen w-screen items-start justify-center xl:items-center">
      <canvas
        id="screen"
        ref={canvasRef}
        className={clsx(
          'block min-h-[50vh] min-w-[50vw] max-w-full select-none object-scale-down',
          mouseStyle
        )}
        style={{ transform: `scale(${videoParameters.scale})` }}
      ></canvas>
    </div>
  );
};
