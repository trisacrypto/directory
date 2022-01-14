import NoDataImg from 'assets/images/no-data-rafiki.svg'
import PropTypes from 'prop-types';

const NoData = ({ title }) => {
    return (
        <div className='text-center'>
            <div>
                <img src={NoDataImg} width={300} alt="no data" />
            </div>
            <p>
                {title}
            </p>
        </div>
    )
}

NoData.propTypes = {
    title: PropTypes.string
}

NoData.defaultProps = {
    title: 'No Data'
}

export default NoData