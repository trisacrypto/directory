import { Meta, Story } from '@storybook/react';
import ConfirmIdentityCertificateModal, { ConfirmIdentityCertificateProps } from '.';

export default {
  title: 'components/ConfirmIdentityCertificateRequest',
  component: ConfirmIdentityCertificateModal
} as Meta<ConfirmIdentityCertificateProps>;

const Template: Story<ConfirmIdentityCertificateProps> = (args) => (
  <ConfirmIdentityCertificateModal {...args} />
);

export const Standard = Template.bind({});
Standard.args = {};
