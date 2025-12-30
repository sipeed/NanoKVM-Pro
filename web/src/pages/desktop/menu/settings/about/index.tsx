import { Divider } from 'antd';

import { Community } from './community.tsx';
import { Header } from './header.tsx';
import { Hostname } from './hostname.tsx';
import { Information } from './information.tsx';

export const About = () => {
  return (
    <>
      <Header />
      <Divider className="opacity-50" />

      <Information />
      <Hostname />
      <Divider className="opacity-50" />

      <Community />
    </>
  );
};
