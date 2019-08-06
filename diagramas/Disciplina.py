# coding=System
from timestamp import *

class Disciplina(object):

  """
   

  :version:
  :author:
  """

  """ ATTRIBUTES

   

  id  (public)

   

  descricao  (public)

   

  iduems  (public)

   

  ano  (public)

   

  createdAt  (public)

   

  updatedAt  (public)

  """

  def Create(self, db):
    """
     

    @param *sql.DB db : 
    @return (int64,error) :
    @author
    """
    func (this *Disciplina) Create(db *sql.DB) (int64, error) {
    	var id int64
    	stmtIns, err := db.Prepare("INSERT INTO cadastros.disciplinas (descricao,id_uems,ano,created_at,updated_at) VALUES ($1,$2,$3,$4,$5) RETURNING id")
    	if err != nil {
    		return id, err
    	}
    	if err := stmtIns.QueryRow(
    		this.Descricao,
    		this.IDUEMS,
    		this.Ano,
    		this.CreatedAt,
    		this.UpdatedAt,
    	).Scan(&id); err != nil {
    		return id, err
    	}
    	defer stmtIns.Close()
    	return id, nil
    }


  def isExist(self, db):
    """
     

    @param *sql.DB db : 
    @return bool :
    @author
    """
    func (this Disciplina) IsExist(db *sql.DB) bool {
    	var count bool
    	_ = db.QueryRow("select exists(select 1 from cadastros.disciplinas where id_uems=$1 and Descricao=$2)", this.IDUEMS, this.Descricao).Scan(&count)
    	return count
    }


  def GetByIDUEMS(self, db):
    """
     

    @param *sql.DB db : 
    @return error :
    @author
    """
    func (this *Disciplina) GetByIDUEMS(id_uems int64, db *sql.DB) error {
    	sql := "SELECT id, descricao, id_uems, ano,created_at,updated_at FROM cadastros.disciplinas where id_uems = $1"
    	selDB, err := db.Query(sql, id_uems)
    	if err != nil {
    		return err
    	}
    	for selDB.Next() {
    		err = selDB.Scan(
    			&this.ID,
    			&this.Descricao,
    			&this.IDUEMS,
    			&this.Ano,
    			&this.CreatedAt,
    			&this.UpdatedAt,
    		)
    		if err != nil {
    			return err
    		}
    	}
    	return nil
    }
    




