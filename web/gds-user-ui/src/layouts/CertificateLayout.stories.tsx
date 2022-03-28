import { Meta, Story } from '@storybook/react';
import { withReduxContext } from 'hoc/withReduxContext';
import CertificateLayout from 'layouts/CertificateLayout';

type CertificateLayoutProps = {
  children: React.ReactNode;
};

export default {
  title: 'layouts/CertificateLayout',
  component: CertificateLayout,
  decorators: [withReduxContext()]
} as Meta;

const Template: Story<CertificateLayoutProps> = (args) => <CertificateLayout {...args} />;

export const Standard = Template.bind({});
Standard.args = {};
