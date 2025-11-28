import { useEffect, useState } from 'react';
import { GithubOutlined } from '@ant-design/icons';
import { Button, Divider, Modal, Tooltip } from 'antd';
import { BotMessageSquare } from 'lucide-react';
import { useTranslation } from 'react-i18next';

import * as api from '@/api/extensions/assistant';

export function AIAssistant() {
  const { t } = useTranslation();
  const [modalOpen, setModalOpen] = useState(false);
  const [content, setContent] = useState<string[]>([]);
  const [intro, setIntro] = useState<string[]>([]);
  const [isLoading, setIsLoading] = useState<'' | 'install' | 'run'>('');

  function install() {
    if (isLoading !== '') return;
    setIsLoading('install');

    api
      .install()
      .then((rsp) => {
        if (rsp.code !== 0) {
          console.log(rsp.msg);
          setIsLoading('');
          return;
        }

        setTimeout(() => {
          window.open(`/#terminal?type=assistant&cmd=install`, '_blank');
          setIsLoading('');
        }, 2000);
      })
      .catch(() => {
        setIsLoading('');
      });
  }

  function start() {
    if (isLoading !== '') return;
    setIsLoading('run');

    api
      .start()
      .then((rsp) => {
        if (rsp.code !== 0) {
          console.log(rsp.msg);
          return;
        }

        window.open(`http://${window.location.host.replace(/:\d+/, '')}:5000/`);
      })
      .catch((err: any) => {
        setIsLoading('');
        console.error('Failed to start AI Assistant server: ', err);
        return;
      })
      .finally(() => {
        setIsLoading('');
      });
  }

  useEffect(() => {
    const contentText = t('assistant.tips.content').split('\n');
    setContent(contentText);
    const introText = t('assistant.tips.intro').split('\n');
    setIntro(introText);
  }, [t]);

  return (
    <div className="flex">
      <Tooltip placement="top" title={t('assistant.title')}>
        <div
          className="flex h-[30px] w-[30px] cursor-pointer items-center justify-center rounded text-neutral-300 hover:bg-neutral-700/80 hover:text-white"
          onClick={() => setModalOpen(true)}
        >
          <BotMessageSquare size={18} />
        </div>
      </Tooltip>

      <Modal width={800} open={modalOpen} onCancel={() => setModalOpen(false)} footer={null}>
        <div className="mt-6 max-h-[650px] overflow-y-auto px-1 md:px-3">
          <div className="text-3xl font-bold">{t('assistant.title')}</div>
          <Divider />
          <div className="content">
            {intro.map((text, index) => (
              <p className="intro indent-8 text-neutral-200/75" key={index}>
                {text}
              </p>
            ))}
          </div>

          <Divider />

          <p className="flex w-full justify-center text-center text-lg text-amber-400">
            ⚠⚠ {t('assistant.tips.warning')} ⚠⚠
          </p>
          {
            <>
              {content.map((t) => (
                <p className="indent-8">{t}</p>
              ))}
            </>
          }
          <div className="flex gap-2">
            <span>{t('assistant.tips.contribute')}</span>
            <a
              className="text-neutral-300/50 hover:text-blue-500"
              href="https://github.com/sipeed/NanoKVM-Pro/tree/main/support/assistant"
              target="_blank"
            >
              <GithubOutlined />
            </a>
          </div>

          <Divider />

          <div className="flex items-center justify-center space-x-5 pb-5 pt-2">
            <Button type="primary" size="large" loading={isLoading === 'install'} onClick={install}>
              {t('assistant.install')}
            </Button>

            <Button type="primary" size="large" loading={isLoading === 'run'} onClick={start}>
              {t('assistant.tryNow')}
            </Button>
          </div>
        </div>
      </Modal>
    </div>
  );
}
