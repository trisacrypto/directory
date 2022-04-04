import { Meta, Story } from '@storybook/react';
import ReviewSubmit from '.';
interface ReviewSubmitProps {
  onSubmitHandler: (e: React.FormEvent, network: string) => void;
}
export default {
  title: 'components/ReviewSubmit',
  component: ReviewSubmit
} as Meta;

const Template: Story<ReviewSubmitProps> = (args) => <ReviewSubmit {...args} />;

export const Standard = Template.bind({});
Standard.args = {};
