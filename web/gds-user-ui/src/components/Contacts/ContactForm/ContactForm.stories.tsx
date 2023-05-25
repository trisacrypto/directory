import { Meta, Story } from '@storybook/react';
import ContactForm from '.';

type ContactFormProps = {
  title: string;
  description: string;
  name: string;
};

export default {
  title: 'components/Contact Form',
  component: ContactForm
} as Meta<ContactFormProps>;

const Template: Story<ContactFormProps> = (args) => <ContactForm {...args} />;

export const Default = Template.bind({});
Default.args = {
  title: 'Legal/Compliance Contact (default)',
  description:
    'Compliance officer or legal contact for requests about the compliance requirements and legal status of your organization.'
};
