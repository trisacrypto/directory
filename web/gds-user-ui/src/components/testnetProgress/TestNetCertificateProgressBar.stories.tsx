import { Meta, Story } from "@storybook/react";

import TestNetCertificateProgressBar from "./TestNetCertificateProgressBar.component";

type TestNetCertificateProgressBarProps = {};

export default {
  title: "components/TestNetCertificateProgressBar",
  component: TestNetCertificateProgressBar,
} as Meta<TestNetCertificateProgressBarProps>;

const Template: Story<TestNetCertificateProgressBarProps> = (args) => (
  <TestNetCertificateProgressBar {...args} />
);

export const Default = Template.bind({});
Default.args = {};
