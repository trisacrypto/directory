import { Meta, Story } from "@storybook/react";
import CertificateManagement from "./index";

export default {
    title: 'components/CertificateManagement',
    component: CertificateManagement
} as Meta;

const Template: Story = args => <CertificateManagement {...args} />;

export const Standard = Template.bind({});
Standard.args = {};
