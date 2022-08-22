import { Meta, Story } from '@storybook/react';

import TrisaImplementationReview from './TrisaImplementationReview';
interface TrisaImplementationReviewProps {
  data: any;
}
export default {
  title: 'components/TrisaImplementationReview',
  component: TrisaImplementationReview
} as Meta;

const Template: Story<TrisaImplementationReviewProps> = (args) => (
  <TrisaImplementationReview {...args} />
);

export const Default = Template.bind({});
Default.args = {};
