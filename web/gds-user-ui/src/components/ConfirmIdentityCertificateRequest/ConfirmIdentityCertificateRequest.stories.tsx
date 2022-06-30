import { Meta, Story } from '@storybook/react';
import ConfirmIdentityCertificate from '.';

export default {
  title: 'components/ConfirmIdentityCertificateRequest',
  component: ConfirmIdentityCertificate
} as Meta;

const Template: Story = (args) => <ConfirmIdentityCertificate {...args} />;

export const Standard = Template.bind({});
Standard.args = {};
