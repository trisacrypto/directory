import { Meta, Story } from "@storybook/react";
import CertificateSection from ".";

export default {
  title: "components/CertificateSection",
  component: CertificateSection,
} as Meta;

const Template: Story = (args) => <CertificateSection {...args} />;

export const Default = Template.bind({});
Default.args = {
  step: 2,
};
