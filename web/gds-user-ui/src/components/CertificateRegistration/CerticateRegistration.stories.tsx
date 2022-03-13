import { Meta, Story } from "@storybook/react";

import CertificateRegistration from ".";

export default {
  title: "components/CertificateRegistration",
  component: CertificateRegistration,
} as Meta;

const Template: Story = (args) => <CertificateRegistration {...args} />;

export const Default = Template.bind({});
Default.args = {};
