import { Meta, Story } from '@storybook/react';
import ConfirmIdentityCertificate, { ConfirmIdentityCertificateProps } from '.';

export default {
  title: 'components/ConfirmIdentityCertificateRequest',
  component: ConfirmIdentityCertificate
} as Meta<ConfirmIdentityCertificateProps>;

const Template: Story<ConfirmIdentityCertificateProps> = (args) => (
  <ConfirmIdentityCertificate {...args} />
);

export const Standard = Template.bind({});
Standard.args = {};
