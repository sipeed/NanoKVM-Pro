import { Divider } from 'antd';

import { Community } from './community.tsx';
import { Header } from './header.tsx';
import { Hostname } from './hostname.tsx';
import { Information } from './information.tsx';

export const About = () => {
  return (
    <>
      <Header />
      <Divider />

      <Information />
      <Hostname />

      <Divider style={{ margin: '32px 0' }} />

      <Community />
    </>
  );
};
