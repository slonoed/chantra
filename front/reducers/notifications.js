import {createAction, handleActions, combineActions} from 'redux-actions'

// ACTIONS

// uniq notification id
let id = 0

const show = createAction('notification_show')
const hide = createAction('notification_hide')
export const showNotification = text => dispatch => {
  const localId = id++
  dispatch(show({text, id: localId}))
  setTimeout(() => dispatch(hide({id: localId})), 3000)
}

// REDUCERS

const initialState = {
  // id -> notification
}

export default handleActions(
  {
    [show]: (state, {payload: {text, id}}) => {
      return {
        ...state,
        [id]: text,
      }
    },
    [hide]: (state, {payload: {id}}) => {
      return {
        ...state,
        [id]: null,
      }
    },
  },
  initialState
)
