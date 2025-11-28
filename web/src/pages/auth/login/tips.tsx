import { useState } from 'react';
import { Button, Card, Modal, Typography } from 'antd';
import { useTranslation } from 'react-i18next';

export const Tips = () => {
  const { t } = useTranslation();
  const [isModalOpen, setIsModalOpen] = useState(false);

  const showModal = () => {
    setIsModalOpen(true);
  };

  const hideModal = () => {
    setIsModalOpen(false);
  };

  return (
    <>
      <span
        className="cursor-pointer text-neutral-300 underline underline-offset-4"
        onClick={showModal}
      >
        {t('auth.forgetPassword')}
      </span>

      <Modal
        title={t('auth.forgetPassword')}
        open={isModalOpen}
        onCancel={hideModal}
        closeIcon={null}
        footer={null}
        centered={true}
      >
        <Card style={{ marginTop: '20px' }}>
          <div className="flex flex-col space-y-3 md:w-[430px]">
            <ul className="list-inside list-disc ps-0">
              <li>
                {t('auth.tips.reset1')}
                <Typography.Text code={true}>admin/admin</Typography.Text>
              </li>
              <li>
                {t('auth.tips.reset2')}
                <Typography.Text code={true}>root/sipeed</Typography.Text>
              </li>
              <li>
                {t('auth.tips.reset3')}
                <Typography.Link
                  href="https://wiki.sipeed.com/hardware/en/kvm/NanoKVM_Pro/extended.html#Factory-Reset"
                  target="_blank"
                >
                  Reset
                </Typography.Link>
              </li>
            </ul>
          </div>
        </Card>

        <div className="flex justify-center pb-3 pt-10">
          <Button type="primary" className="w-24" onClick={hideModal}>
            {t('auth.ok')}
          </Button>
        </div>
      </Modal>
    </>
  );
};
