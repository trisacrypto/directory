import { Meta, Story } from '@storybook/react';
import { withReduxContext } from 'hoc/withReduxContext';

import TestNetCertificateProgressBar from './TestNetCertificateProgressBar.component';

type TestNetCertificateProgressBarProps = {};

export default {
  title: 'components/TestNetCertificateProgressBar',
  component: TestNetCertificateProgressBar,
  decorators: [withReduxContext()]
} as Meta<TestNetCertificateProgressBarProps>;

const Template: Story<TestNetCertificateProgressBarProps> = (args) => (
  <TestNetCertificateProgressBar {...args} />
);

export const Default = Template.bind({});
Default.args = {};
