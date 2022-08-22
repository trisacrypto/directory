import { Meta, Story } from '@storybook/react';

import ContactsReview from './ContactsReview';

interface ContactsReviewProps {
  data: any;
}
export default {
  title: 'components/ContactsReview',
  component: ContactsReview
} as Meta;

const Template: Story<ContactsReviewProps> = (args) => <ContactsReview {...args} />;

export const Default = Template.bind({});
Default.args = {};
