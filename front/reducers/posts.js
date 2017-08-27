import {createAction, handleActions, combineActions} from 'redux-actions'
import {showNotification} from './notifications'

// ACTIONS

const scheduleNewPostStart = createAction('scheduleNewPostStart')
const scheduleNewPostSuccess = createAction('scheduleNewPostSuccess')
const scheduleNewPostFail = createAction('scheduleNewPostFail')
export const scheduleNewPost = (post, channelId) => dispatch => {
  dispatch(scheduleNewPostStart())

  const content = Object.assign(post, {
    sentDate: post.sentDate.unix(),
    channel_id: channelId,
  })

  fetch('/api/addPost', {
    method: 'POST',
    body: JSON.stringify(content),
    credentials: 'include',
  })
    .then(r => r.json())
    .then(r => {
      if (!r.ok) {
        throw new Error(r.error || 'Unknown error')
      } else {
        dispatch(scheduleNewPostSuccess({post: r.post}))
        dispatch(showNotification('Created!'))
      }
    })
    .catch(error => {
      dispatch(scheduleNewPostFail({error}))
    })
}

// REDUCERS

const initialState = {
  // id -> post
  posts: (INITIAL_DATA.posts || []).reduce((all, c) => {
    all[c.id] = c
    return all
  }, {}),

  form: {
    loading: false,
  },
}

export default handleActions(
  {
    [scheduleNewPostStart]: (state, action) => {
      return {
        ...state,
        form: {
          ...state.form,
          loading: true,
        },
      }
    },
    [scheduleNewPostSuccess]: (state, {payload}) => {
      const {post} = payload

      return {
        ...state,
        form: {
          ...state.form,
          loading: false,
        },
        posts: {
          ...state.posts,
          [post.id]: post,
        },
      }
    },
    [scheduleNewPostFail]: (state, {payload: {error}}) => {
      return {
        ...state,
        form: {
          ...state.form,
          loading: false,
          error: error ? error.toString() : 'Unknown error',
        },
      }
    },
  },
  initialState
)
