//!  misbehaviour handling.

use crate::{
    error::Error,
    query::check_for_misbehaviour,
    state::{ClientState, Header},
};

/// Checks two headers and freezes only a valid same-height conflict.
pub fn apply<Config>(
    client: &mut ClientState<Config>,
    first: &Header,
    second: &Header,
) -> Result<bool, Error> {
    if !check_for_misbehaviour(first, second)? {
        return Ok(false);
    }
    client.frozen_height = Some(first.height.revision_height);
    Ok(true)
}
