
/* 
 * for saving, restoring a player's game session. 
 * a global application active count increases monotonically 
 * for every request to activate any domain.
 * a particular domain is in scope when its own active value is non-zero; 
 * its value gets reset to zero when the domain falls out of scope.
 */
create table if not exists 
	rt.run_domain( domain text, active int, primary key( domain )); 

/* we dont need an "active" -- we can join against run_domain, or write 0 to domain to disable a pair. */
create table if not exists 
	rt.run_pair( domain text, relKind int, oneNoun int, otherNoun int, unique( relKind, oneNoun, otherNoun ) ); 

/* fields a s*/
create table if not exists 
	rt.run_value( domain text, noun text, field text, value blob );