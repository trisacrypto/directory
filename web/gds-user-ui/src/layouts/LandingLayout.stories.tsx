import React from 'react';
import { Story } from '@storybook/react';

import LandingLayout from './LandingLayout';

type LandingLayoutProps = {
  children: React.ReactNode;
};

export default {
  title: 'Layouts/LandingLayout',
  component: LandingLayout
};

export const Default: Story<LandingLayoutProps> = (props) => <LandingLayout {...props} />;

Default.bind({
  children: <div>Hello World</div>
});
